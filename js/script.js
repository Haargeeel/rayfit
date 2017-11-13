window.onload = function() {
  const orbitISS = 382.5
  const orbitHubble = 595
  const maxDist = 600

  const distanceDiv = document.getElementById('distance')
  const me = document.getElementById('me')
  const iss = document.getElementById('iss')
  const issWrapper = document.getElementById('issWrapper')
  const hubbleWrapper = document.getElementById('hubbleWrapper')
  const hubble = document.getElementById('hubble')

  const distanceWidth = distanceDiv.getBoundingClientRect().width

  me.style.left = `${myDistance * distanceWidth / maxDist}px`
  me.style.display = 'block'

  issWrapper.style.width = `${2 * orbitISS * distanceWidth / maxDist + 600}px`
  issWrapper.style.height = `${2 * orbitISS * distanceWidth / maxDist + 600}px`
  // iss.style.left = `${orbitISS * distanceWidth / orbitHubble + 200}px`
  iss.style.display = 'block'

  hubbleWrapper.style.width = `${2 * orbitHubble * distanceWidth / maxDist + 600}px`
  hubbleWrapper.style.height = `${2 * orbitHubble * distanceWidth / maxDist + 600}px`
  // hubble.style.left = `${orbitHubble * distanceWidth / maxDist}px`
  hubble.style.display = 'block'
}
