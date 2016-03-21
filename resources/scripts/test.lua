

function update (time, ctx)
    local query           = ctx:EntityManager():MakeCmptQuery({'physics', 'position'})
    local positionCmptSet = ctx:EntityManager():GetComponentSet('position')

    local eids = ctx:EntityManager():EntitiesMatching(query)
    local idxs = positionCmptSet:Indices(eids)

    local positionCmptSlice = positionCmptSet:Slice()
    local cmpt = positionCmptSlice[idxs[1]]
    cmpt:SetPos({2, 2})
    -- cmpt.Pos[1] = cmpt.Pos[1] + 0.01
    print('pos = ' .. tostring(cmpt.Pos))


    -- for key, val in ipairs(idxs) do
    --     print('key: ' .. tostring(key) .. ' // val: ' .. tostring(val))
    -- end
end