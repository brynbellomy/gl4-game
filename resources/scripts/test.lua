

function update (time, ctx, ownID)
    local em = ctx:EntityManager()

    local tagsys = em:GetSystem('tag')
    local heroID = tagsys:EntityWithTag('hero')
    local posCmptSet = em:GetComponentSet('position')
    local moveCmptSet = em:GetComponentSet('move')

    local moveIdx = moveCmptSet:Index(ownID)
    local moveCmptSlice = moveCmptSet:Slice()
    moveCmptSlice[moveIdx]:SetVector({-1, -1})

    local heroPosCmpt, err = posCmptSet:Get(heroID)
    local selfPosCmpt, err = posCmptSet:Get(ownID)
    local selfMoveCmpt, err = moveCmptSet:Get(ownID)
    -- print('heroCmpt = ' .. tostring(heroPosCmpt))
    -- selfMoveCmpt:SetVector({-1, -1})
    -- selfMoveCmpt.Vec = {-1, -1}




    -- local query           = ctx:EntityManager():MakeCmptQuery({'physics', 'position'})
    -- local posCmptSet = ctx:EntityManager():GetComponentSet('position')

    -- local eids = ctx:EntityManager():EntitiesMatching(query)
    -- local idxs = posCmptSet:Indices(eids)

    -- local positionCmptSlice = posCmptSet:Slice()
    -- local cmpt = positionCmptSlice[idxs[1]]
    -- cmpt:SetPos({2, 2})
    -- cmpt.Pos[1] = cmpt.Pos[1] + 0.01
    -- print('pos = ' .. tostring(cmpt.Pos))


    -- for key, val in ipairs(idxs) do
    --     print('key: ' .. tostring(key) .. ' // val: ' .. tostring(val))
    -- end
end