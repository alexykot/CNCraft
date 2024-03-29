// Code generated by "stringer -type=PacketType packets.go"; DO NOT EDIT.

package protocol

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TypeUnspecified - -1]
	_ = x[SHandshake-4096]
	_ = x[SRequest-4352]
	_ = x[SPing-4353]
	_ = x[SLoginStart-4608]
	_ = x[SEncryptionResponse-4609]
	_ = x[SLoginPluginResponse-4610]
	_ = x[STeleportConfirm-4864]
	_ = x[SQueryBlockNBT-4865]
	_ = x[SQueryEntityNBT-4877]
	_ = x[SSetDifficulty-4866]
	_ = x[SChatMessage-4867]
	_ = x[SClientStatus-4868]
	_ = x[SClientSettings-4869]
	_ = x[STabComplete-4870]
	_ = x[SWindowConfirmation-4871]
	_ = x[SClickWindowButton-4872]
	_ = x[SClickWindow-4873]
	_ = x[SCloseWindow-4874]
	_ = x[SPluginMessage-4875]
	_ = x[SEditBook-4876]
	_ = x[SInteractEntity-4878]
	_ = x[SGenerateStructure-4879]
	_ = x[SKeepAlive-4880]
	_ = x[SLockDifficulty-4881]
	_ = x[SPlayerPosition-4882]
	_ = x[SPlayerPosAndRotation-4883]
	_ = x[SPlayerRotation-4884]
	_ = x[SPlayerMovement-4885]
	_ = x[SVehicleMove-4886]
	_ = x[SSteerBoat-4887]
	_ = x[SPickItem-4888]
	_ = x[SCraftRecipeRequest-4889]
	_ = x[SPlayerAbilities-4890]
	_ = x[SPlayerDigging-4891]
	_ = x[SEntityAction-4892]
	_ = x[SSteerVehicle-4893]
	_ = x[SSetDisplayedRecipe-4894]
	_ = x[SSetRecipeBookState-4895]
	_ = x[SNameItem-4896]
	_ = x[SResourcePackStatus-4897]
	_ = x[SAdvancementTab-4898]
	_ = x[SSelectTrade-4899]
	_ = x[SSetBeaconEffect-4900]
	_ = x[SHeldItemChange-4901]
	_ = x[SUpdateCommandBlock-4902]
	_ = x[SUpdateCommandBlockMinecart-4903]
	_ = x[SCreativeInventoryAction-4904]
	_ = x[SUpdateJigsawBlock-4905]
	_ = x[SUpdateStructureBlock-4906]
	_ = x[SUpdateSign-4907]
	_ = x[SAnimation-4908]
	_ = x[SSpectate-4909]
	_ = x[SPlayerBlockPlacement-4910]
	_ = x[SUseItem-4911]
	_ = x[CResponse-61696]
	_ = x[CPong-61697]
	_ = x[CDisconnectLogin-61952]
	_ = x[CEncryptionRequest-61953]
	_ = x[CLoginSuccess-61954]
	_ = x[CSetCompression-61955]
	_ = x[CLoginPluginRequest-61956]
	_ = x[CSpawnEntity-62208]
	_ = x[CSpawnExperienceOrb-62209]
	_ = x[CSpawnLivingEntity-62210]
	_ = x[CSpawnPainting-62211]
	_ = x[CSpawnPlayer-62212]
	_ = x[CEntityAnimation-62213]
	_ = x[CStatistics-62214]
	_ = x[CAcknowledgePlayerDigging-62215]
	_ = x[CBlockBreakAnimation-62216]
	_ = x[CBlockEntityData-62217]
	_ = x[CBlockAction-62218]
	_ = x[CBlockChange-62219]
	_ = x[CBossBar-62220]
	_ = x[CServerDifficulty-62221]
	_ = x[CChatMessage-62222]
	_ = x[CTabComplete-62223]
	_ = x[CDeclareCommands-62224]
	_ = x[CWindowConfirmation-62225]
	_ = x[CCloseWindow-62226]
	_ = x[CWindowItems-62227]
	_ = x[CWindowProperty-62228]
	_ = x[CSetSlot-62229]
	_ = x[CSetCooldown-62230]
	_ = x[CPluginMessage-62231]
	_ = x[CNamedSoundEffect-62232]
	_ = x[CDisconnectPlay-62233]
	_ = x[CEntityStatus-62234]
	_ = x[CExplosion-62235]
	_ = x[CUnloadChunk-62236]
	_ = x[CChangeGameState-62237]
	_ = x[COpenHorseWindow-62238]
	_ = x[CKeepAlive-62239]
	_ = x[CChunkData-62240]
	_ = x[CEffect-62241]
	_ = x[CParticle-62242]
	_ = x[CUpdateLight-62243]
	_ = x[CJoinGame-62244]
	_ = x[CMapData-62245]
	_ = x[CTradeList-62246]
	_ = x[CEntityPosition-62247]
	_ = x[CEntityPositionandRotation-62248]
	_ = x[CEntityRotation-62249]
	_ = x[CEntityMovement-62250]
	_ = x[CVehicleMove-62251]
	_ = x[COpenBook-62252]
	_ = x[COpenWindow-62253]
	_ = x[COpenSignEditor-62254]
	_ = x[CCraftRecipeResponse-62255]
	_ = x[CPlayerAbilities-62256]
	_ = x[CCombatEvent-62257]
	_ = x[CPlayerInfo-62258]
	_ = x[CFacePlayer-62259]
	_ = x[CPlayerPositionAndLook-62260]
	_ = x[CUnlockRecipes-62261]
	_ = x[CDestroyEntities-62262]
	_ = x[CRemoveEntityEffect-62263]
	_ = x[CResourcePackSend-62264]
	_ = x[CRespawn-62265]
	_ = x[CEntityHeadLook-62266]
	_ = x[CMultiBlockChange-62267]
	_ = x[CSelectAdvancementTab-62268]
	_ = x[CWorldBorder-62269]
	_ = x[CCamera-62270]
	_ = x[CHeldItemChange-62271]
	_ = x[CUpdateViewPosition-62272]
	_ = x[CUpdateViewDistance-62273]
	_ = x[CSpawnPosition-62274]
	_ = x[CDisplayScoreboard-62275]
	_ = x[CEntityMetadata-62276]
	_ = x[CAttachEntity-62277]
	_ = x[CEntityVelocity-62278]
	_ = x[CEntityEquipment-62279]
	_ = x[CSetExperience-62280]
	_ = x[CUpdateHealth-62281]
	_ = x[CScoreboardObjective-62282]
	_ = x[CSetPassengers-62283]
	_ = x[CTeams-62284]
	_ = x[CUpdateScore-62285]
	_ = x[CTimeUpdate-62286]
	_ = x[CTitle-62287]
	_ = x[CEntitySoundEffect-62288]
	_ = x[CSoundEffect-62289]
	_ = x[CStopSound-62290]
	_ = x[CPlayerListHeaderAndFooter-62291]
	_ = x[CNBTQueryResponse-62292]
	_ = x[CCollectItem-62293]
	_ = x[CEntityTeleport-62294]
	_ = x[CAdvancements-62295]
	_ = x[CEntityProperties-62296]
	_ = x[CEntityEffect-62297]
	_ = x[CDeclareRecipes-62298]
	_ = x[CTags-62299]
}

const (
	_PacketType_name_0 = "TypeUnspecified"
	_PacketType_name_1 = "SHandshake"
	_PacketType_name_2 = "SRequestSPing"
	_PacketType_name_3 = "SLoginStartSEncryptionResponseSLoginPluginResponse"
	_PacketType_name_4 = "STeleportConfirmSQueryBlockNBTSSetDifficultySChatMessageSClientStatusSClientSettingsSTabCompleteSWindowConfirmationSClickWindowButtonSClickWindowSCloseWindowSPluginMessageSEditBookSQueryEntityNBTSInteractEntitySGenerateStructureSKeepAliveSLockDifficultySPlayerPositionSPlayerPosAndRotationSPlayerRotationSPlayerMovementSVehicleMoveSSteerBoatSPickItemSCraftRecipeRequestSPlayerAbilitiesSPlayerDiggingSEntityActionSSteerVehicleSSetDisplayedRecipeSSetRecipeBookStateSNameItemSResourcePackStatusSAdvancementTabSSelectTradeSSetBeaconEffectSHeldItemChangeSUpdateCommandBlockSUpdateCommandBlockMinecartSCreativeInventoryActionSUpdateJigsawBlockSUpdateStructureBlockSUpdateSignSAnimationSSpectateSPlayerBlockPlacementSUseItem"
	_PacketType_name_5 = "CResponseCPong"
	_PacketType_name_6 = "CDisconnectLoginCEncryptionRequestCLoginSuccessCSetCompressionCLoginPluginRequest"
	_PacketType_name_7 = "CSpawnEntityCSpawnExperienceOrbCSpawnLivingEntityCSpawnPaintingCSpawnPlayerCEntityAnimationCStatisticsCAcknowledgePlayerDiggingCBlockBreakAnimationCBlockEntityDataCBlockActionCBlockChangeCBossBarCServerDifficultyCChatMessageCTabCompleteCDeclareCommandsCWindowConfirmationCCloseWindowCWindowItemsCWindowPropertyCSetSlotCSetCooldownCPluginMessageCNamedSoundEffectCDisconnectPlayCEntityStatusCExplosionCUnloadChunkCChangeGameStateCOpenHorseWindowCKeepAliveCChunkDataCEffectCParticleCUpdateLightCJoinGameCMapDataCTradeListCEntityPositionCEntityPositionandRotationCEntityRotationCEntityMovementCVehicleMoveCOpenBookCOpenWindowCOpenSignEditorCCraftRecipeResponseCPlayerAbilitiesCCombatEventCPlayerInfoCFacePlayerCPlayerPositionAndLookCUnlockRecipesCDestroyEntitiesCRemoveEntityEffectCResourcePackSendCRespawnCEntityHeadLookCMultiBlockChangeCSelectAdvancementTabCWorldBorderCCameraCHeldItemChangeCUpdateViewPositionCUpdateViewDistanceCSpawnPositionCDisplayScoreboardCEntityMetadataCAttachEntityCEntityVelocityCEntityEquipmentCSetExperienceCUpdateHealthCScoreboardObjectiveCSetPassengersCTeamsCUpdateScoreCTimeUpdateCTitleCEntitySoundEffectCSoundEffectCStopSoundCPlayerListHeaderAndFooterCNBTQueryResponseCCollectItemCEntityTeleportCAdvancementsCEntityPropertiesCEntityEffectCDeclareRecipesCTags"
)

var (
	_PacketType_index_2 = [...]uint8{0, 8, 13}
	_PacketType_index_3 = [...]uint8{0, 11, 30, 50}
	_PacketType_index_4 = [...]uint16{0, 16, 30, 44, 56, 69, 84, 96, 115, 133, 145, 157, 171, 180, 195, 210, 228, 238, 253, 268, 289, 304, 319, 331, 341, 350, 369, 385, 399, 412, 425, 444, 463, 472, 491, 506, 518, 534, 549, 568, 595, 619, 637, 658, 669, 679, 688, 709, 717}
	_PacketType_index_5 = [...]uint8{0, 9, 14}
	_PacketType_index_6 = [...]uint8{0, 16, 34, 47, 62, 81}
	_PacketType_index_7 = [...]uint16{0, 12, 31, 49, 63, 75, 91, 102, 127, 147, 163, 175, 187, 195, 212, 224, 236, 252, 271, 283, 295, 310, 318, 330, 344, 361, 376, 389, 399, 411, 427, 443, 453, 463, 470, 479, 491, 500, 508, 518, 533, 559, 574, 589, 601, 610, 621, 636, 656, 672, 684, 695, 706, 728, 742, 758, 777, 794, 802, 817, 834, 855, 867, 874, 889, 908, 927, 941, 959, 974, 987, 1002, 1018, 1032, 1045, 1065, 1079, 1085, 1097, 1108, 1114, 1132, 1144, 1154, 1180, 1197, 1209, 1224, 1237, 1254, 1267, 1282, 1287}
)

func (i PacketType) String() string {
	switch {
	case i == -1:
		return _PacketType_name_0
	case i == 4096:
		return _PacketType_name_1
	case 4352 <= i && i <= 4353:
		i -= 4352
		return _PacketType_name_2[_PacketType_index_2[i]:_PacketType_index_2[i+1]]
	case 4608 <= i && i <= 4610:
		i -= 4608
		return _PacketType_name_3[_PacketType_index_3[i]:_PacketType_index_3[i+1]]
	case 4864 <= i && i <= 4911:
		i -= 4864
		return _PacketType_name_4[_PacketType_index_4[i]:_PacketType_index_4[i+1]]
	case 61696 <= i && i <= 61697:
		i -= 61696
		return _PacketType_name_5[_PacketType_index_5[i]:_PacketType_index_5[i+1]]
	case 61952 <= i && i <= 61956:
		i -= 61952
		return _PacketType_name_6[_PacketType_index_6[i]:_PacketType_index_6[i+1]]
	case 62208 <= i && i <= 62299:
		i -= 62208
		return _PacketType_name_7[_PacketType_index_7[i]:_PacketType_index_7[i+1]]
	default:
		return "PacketType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
