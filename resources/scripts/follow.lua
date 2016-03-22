

function update (time, ctx, ownID)
    local em = ctx:EntityManager()

    local heroID = em:GetSystem('tag'):EntityWithTag('hero')
    local posCmptSet = em:GetComponentSet('position')
    local moveCmptSet = em:GetComponentSet('move')

    local heroPosIdx = posCmptSet:Index(heroID)
    local ownPosIdx = posCmptSet:Index(ownID)
    local posCmptSlice = posCmptSet:Slice()

    local vec = posCmptSlice[heroPosIdx]:GetPos():Sub( posCmptSlice[ownPosIdx]:GetPos() ):Normalize()
    print('vec = ' .. tostring(vec))

    local ownMoveIdx = moveCmptSet:Index(ownID)
    local moveCmptSlice = moveCmptSet:Slice()
    moveCmptSlice[ownMoveIdx]:SetVector(vec)
end