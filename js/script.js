window.onload = function() {
  const orbitISS = 382.5
  const distanceDiv = document.getElementById('distance')
  const me = document.getElementById('me')
  const distanceWidth = distanceDiv.getBoundingClientRect().width
  me.style.left = `${myDistance * distanceWidth / orbitISS}px`
  me.style.display = 'block'
}
