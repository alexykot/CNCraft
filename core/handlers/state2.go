package handlers

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexykot/cncraft/core/nats"
	"github.com/alexykot/cncraft/core/nats/subj"
	"github.com/alexykot/cncraft/pkg/envelope"
	"github.com/alexykot/cncraft/pkg/envelope/pb"
	"github.com/alexykot/cncraft/pkg/protocol"
	"github.com/alexykot/cncraft/pkg/protocol/auth"
)

// HandleSLoginStart handles the LoginStart packet.
func HandleSLoginStart(auther auth.A, connID uuid.UUID, sPacket protocol.SPacket) ([]protocol.CPacket, error) {
	loginStart, ok := sPacket.(*protocol.SPacketLoginStart)
	if !ok {
		return nil, fmt.Errorf("received packet is not a loginStart: %v", sPacket)
	}

	userID := connID // By design connection ID is also the auth user ID and then the player ID.
	if err := auther.BootstrapUser(userID, loginStart.Username); err != nil {
		return nil, fmt.Errorf("failed to bootstrap user: %w", err)
	}

	if currentConf.IsCracked { // "cracked" or "offline-mode" server does not do authentication or encryption
		loginSuccess, _ := protocol.GetPacketFactory().MakeCPacket(protocol.CLoginSuccess) // Predefined packet is expected to always exist.
		loginSuccess.(*protocol.CPacketLoginSuccess).PlayerUUID = userID.String()          // And always be of the correct type.
		loginSuccess.(*protocol.CPacketLoginSuccess).PlayerName = loginStart.Username
		return []protocol.CPacket{loginSuccess}, nil
	}

	cpacket, _ := protocol.GetPacketFactory().MakeCPacket(protocol.CEncryptionRequest) // Predefined packet is expected to always exist.
	encRequest := cpacket.(*protocol.CPacketEncryptionRequest)                         // And always be of the correct type.

	encRequest.ServerID = currentConf.ServerID
	encRequest.PublicKey = auther.GetUserPubkey(userID)
	encRequest.VerifyToken = auther.GetUserVerifyToken(userID)

	return []protocol.CPacket{encRequest}, nil
}

func HandleSEncryptionResponse(auther auth.A, ps nats.PubSub,
	stateSetter func(state protocol.State), encSetter func([]byte) error, compSetter func(),
	connID uuid.UUID, sPacket protocol.SPacket) ([]protocol.CPacket, error) {

	encResponse, ok := sPacket.(*protocol.SPacketEncryptionResponse)
	if !ok {
		return nil, fmt.Errorf("received packet is not a loginStart: %v", sPacket)
	}

	userID := connID // By design connection ID is also the auth user ID and then the player ID.
	if bytes.Compare(encResponse.VerifyToken, auther.GetUserVerifyToken(userID)) != 0 {
		auther.LoginFailure(userID)
		return nil, newPacketError(InvalidLoginErr, errors.New("supplied verify token does not match the saved one"))
	}

	sharedSecret, err := auther.DecryptUserSharedSecret(userID, encResponse.SharedSecret)
	if err != nil {
		auther.LoginFailure(userID)
		return nil, newPacketError(InvalidLoginErr, fmt.Errorf("failed to decrypt user shared secret: %w", err))
	}

	mojangData, err := auther.RunMojangSessionAuth(userID, sharedSecret)
	if err != nil {
		auther.LoginFailure(userID)
		return nil, newPacketError(InvalidLoginErr, fmt.Errorf("failed to run Mojang session server auth: %w", err))
	}

	compSetter()
	if err := encSetter(sharedSecret); err != nil {
		return nil, fmt.Errorf("failed to enable conn encryption: %w", err)
	}

	setCompression, _ := protocol.GetPacketFactory().MakeCPacket(protocol.CSetCompression)       // Predefined packet is expected to always exist.
	setCompression.(*protocol.CPacketSetCompression).Threshold = currentConf.Network.ZipTreshold // And always be of the correct type.

	loginSuccess, _ := protocol.GetPacketFactory().MakeCPacket(protocol.CLoginSuccess)
	loginSuccess.(*protocol.CPacketLoginSuccess).PlayerUUID = mojangData.ProfileID.String()
	loginSuccess.(*protocol.CPacketLoginSuccess).PlayerName = mojangData.Username

	stateSetter(protocol.Play)
	lope := envelope.PlayerLoading(&pb.PlayerLoading{
		Id:        userID.String(),
		ProfileId: mojangData.ProfileID.String(),
		Username:  mojangData.Username,
		// TODO also publish skin data
	})
	if err := ps.Publish(subj.MkPlayerLoading(), lope); err != nil {
		return nil, fmt.Errorf("failed to publish player loading envelope: %w", err)
	}

	auther.LoginSuccess(userID)
	return []protocol.CPacket{setCompression, loginSuccess}, nil
}
