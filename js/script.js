window.onload = function() {
  const orbitISS = 382.5
  const orbitHubble = 595

  const distanceDiv = document.getElementById('distance')
  const me = document.getElementById('me')
  const iss = document.getElementById('iss')
  const hubble = document.getElementById('hubble')

  const distanceWidth = distanceDiv.getBoundingClientRect().width

  me.style.left = `${myDistance * distanceWidth / orbitHubble}px`
  me.style.display = 'block'

  iss.style.left = `${orbitISS * distanceWidth / orbitHubble}px`
  iss.style.display = 'block'

  hubble.style.right = '0px'
  hubble.style.display = 'block'
}
