
function update(time, scriptInterface) {
    var entityManager = scriptInterface.EntityManager()
    var componentQuery = entityManager.MakeCmptQuery(['physics', 'position'])
    var physicsCmptSet = entityManager.GetComponentSet('physics')
    var positionCmptSet = entityManager.GetComponentSet('position')

    var matchIDs = entityManager.EntitiesMatching(componentQuery)
    var physicsIdxs = physicsCmptSet.Indices(matchIDs)
    var positionIdxs = positionCmptSet.Indices(matchIDs)

    var physicsCmptSlice = physicsCmptSet.Slice()
    var positionCmptSlice = positionCmptSet.Slice()

    for (var i = 0; i < physicsIdxs.length; i++) {
        console.log('velocity ~>', physicsCmptSlice[physicsIdxs[i]].Velocity)
        console.log('position ~>', positionCmptSlice[positionIdxs[i]].Pos)
        physicsCmptSlice.SetVelocity(physicsIdxs[i], [0.1, 0.1])
    }
}