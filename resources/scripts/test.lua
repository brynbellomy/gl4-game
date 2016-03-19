

function update (time, ctx)
    local query          = ctx:entityManager():makeCmptQuery({'physics', 'position'})
    local physicsCmptSet = ctx:entityManager():getComponentSet('physics')

    local eids = ctx:entityManager():entitiesMatching(query)
    local idxs = physicsCmptSet:indices(eids)
    for key, val in ipairs(idxs) do
        print('key: ' .. tostring(key) .. ' // val: ' .. tostring(val))
    end
end