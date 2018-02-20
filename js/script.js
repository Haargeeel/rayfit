window.onload = function() {
  const orbitISS = 382.5
  const orbitHubble = 595

  const distanceDiv = document.getElementById('distance')
  const me = document.getElementById('me')
  const iss = document.getElementById('iss')
  const hubble = document.getElementById('hubble')
  const astronaut = me.getElementsByTagName('img')[0]

  const distanceWidth = distanceDiv.getBoundingClientRect().width

  const rounds = Math.floor(myDistance / orbitISS)
  const currentDistance = myDistance % orbitISS

  console.log('rounds', rounds)
  console.log('distance', currentDistance)
  astronaut.className = rounds % 2 !== 0 ? 'flip' : ''

  me.style.left = `${currentDistance * distanceWidth / orbitHubble}px`
  me.style.display = 'block'

  iss.style.left = `${orbitISS * distanceWidth / orbitHubble}px`
  iss.style.display = 'block'

  hubble.style.right = '0px'
  hubble.style.display = 'block'
}
