

local followTag = nil

function init(ctx, params)
    followTag = params.followTag
    print('follow tag = ' .. tostring(followTag))
end

function update (time, ctx, ownID)
    local em = ctx:EntityManager()

    local tgtID = em:GetSystem('tag'):EntityWithTag(followTag)
    local posCmptSet = em:GetComponentSet('position')
    local moveCmptSet = em:GetComponentSet('move')

    local tgtPosIdx = posCmptSet:Index(tgtID)
    local ownPosIdx = posCmptSet:Index(ownID)
    local posCmptSlice = posCmptSet:Slice()

    local vec = posCmptSlice[tgtPosIdx]:GetPos():Sub( posCmptSlice[ownPosIdx]:GetPos() ):Normalize()

    local ownMoveIdx = moveCmptSet:Index(ownID)
    local moveCmptSlice = moveCmptSet:Slice()
    moveCmptSlice[ownMoveIdx]:SetVector(vec)
    moveCmptSlice[ownMoveIdx]:SetMovementType(1)

    local physCmptSet = em:GetComponentSet('physics')
end