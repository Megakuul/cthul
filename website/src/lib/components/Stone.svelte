<script>
  let { ...restProps } = $props();
</script>

<svg viewBox="0 0 2000 2000" fill="none" xmlns="http://www.w3.org/2000/svg" {...restProps}>
  <g id="granit">
    <path fill-rule="evenodd" clip-rule="evenodd" d="M1015.91 765.942L1041.25 452.36L721.337 351L483.775 465.03L552.934 536.311L505.857 487.847L414 595.542L423.593 750.104L452.1 845.129L505.948 968.662L613.642 974.997L1000.08 997.169L1015.91 765.942Z" fill="#042F0B"/>
    <path d="M1000.08 997.169L695.997 683.587L715.002 883.139L1000.08 997.169Z" fill="#9D9D9D" fill-opacity="0.3"/>
    <path d="M715.002 883.139L423.593 750.104L452.1 845.129L505.948 968.662L613.642 974.997L715.002 883.139Z" fill="#9D9D9D" fill-opacity="0.26"/>
    <path d="M1000.08 997.169L715.002 883.139L613.642 974.997L1000.08 997.169Z" fill="#CDCDCD" fill-opacity="0.58"/>
    <path d="M414 595.542L423.593 750.105L715.002 883.139L695.997 683.587L505.857 487.847L414 595.542Z" fill="#C4C4C4" fill-opacity="0.81"/>
    <path d="M1041.25 452.36L1015.91 765.942L1000.08 997.169L483.775 465.03L721.337 351L1041.25 452.36Z" fill="#D7D7D7" fill-opacity="0.31"/>  
  </g>
  <g id="proton">
    <path fill-rule="evenodd" clip-rule="evenodd" d="M1582.9 968.662L1586.06 446.025L1565.74 461.973L1034.92 531.548L1015.91 765.942L1000.08 997.169L1582.9 968.662ZM1016.01 765.894L1496.83 515.994L1161.62 689.923L1161.62 689.923L1016.01 765.894Z" fill="#57056C"/>
    <path d="M1582.9 968.662L1586.06 446.025L1085.6 838.794L1000.08 997.169L1582.9 968.662Z" fill="#E4E4E4" fill-opacity="0.22"/>
    <path d="M1085.6 838.795L1497.37 515.71L1015.91 765.942L1085.6 838.795Z" fill="#9D9D9D" fill-opacity="0.3"/>
    <path d="M1034.92 531.548L1015.91 765.942L1161.62 689.923L1034.92 531.548Z" fill="#C4C4C4" fill-opacity="0.81"/>
    <path d="M1034.92 531.548L1566.58 461.863L1497.37 515.71L1161.62 689.923L1034.92 531.548Z" fill="#9D9D9D" fill-opacity="0.26"/>
    <path d="M1000.08 997.169L1085.6 838.795L1015.91 765.942L1000.08 997.169Z" fill="#E4E4E4" fill-opacity="0.49"/>  
  </g>
  <g id="wave">
    <path fill-rule="evenodd" clip-rule="evenodd" d="M1807.28 1243.59L1589.23 968.662L1323.16 981.332L1000.08 997.17L1060.26 1510.3L1823.63 1643.34L1947.16 1576.82L1931.45 1567.21L1982 1526.14L1829.96 1250.57L1807.28 1243.59Z" fill="#0D65A4"/>
    <path d="M1418.19 1123.87L1000.08 997.17L1310.49 1187.22L1418.19 1123.87Z" fill="#9D9D9D" fill-opacity="0.3"/>
    <path d="M1807.79 1244.23L1589.23 968.662L1323.16 981.332L1418.19 1123.87L1807.79 1244.23Z" fill="#9D9D9D" fill-opacity="0.26"/>
    <path d="M1418.19 1123.87L1323.16 981.332L1000.08 997.17L1418.19 1123.87Z" fill="#CDCDCD" fill-opacity="0.58"/>
    <path d="M1310.49 1187.22L1931.32 1567.32L1982 1526.14L1829.96 1250.57L1418.19 1123.87L1310.49 1187.22Z" fill="#C4C4C4" fill-opacity="0.81"/>
    <path d="M1060.26 1510.3L1000.08 997.17L1947.16 1576.82L1823.63 1643.34L1060.26 1510.3Z" fill="#F0F0F0" fill-opacity="0.26"/>  
  </g>
  <g id="rune">
    <path fill-rule="evenodd" clip-rule="evenodd" d="M417.257 1415.28L445.765 965.494L1000.08 997.169L999.794 997.372L1000.08 997.169L1038.09 1320.25L1069.76 1589.49L686.607 1222.17L920.889 1446.95L626.312 1484.96L483.775 1367.77L635.877 1258.44L417.257 1415.28Z" fill="#5E4A11"/>
    <path d="M686.495 1222.06L483.775 1367.77L626.312 1484.96L920.889 1446.95L686.495 1222.06Z" fill="#C4C4C4" fill-opacity="0.81"/>
    <path d="M1000.08 997.169L787.855 1149.21L1038.09 1320.25L1000.08 997.169Z" fill="#9D9D9D" fill-opacity="0.3"/>
    <path d="M686.495 1222.06L1069.76 1589.49L1038.09 1320.25L686.495 1222.06Z" fill="#9D9D9D" fill-opacity="0.26"/>
    <path d="M1038.09 1320.25L787.855 1149.21L686.495 1222.06L1038.09 1320.25Z" fill="#CDCDCD" fill-opacity="0.58"/>
    <path d="M445.765 965.494L417.257 1415.28L1000.08 997.169L445.765 965.494Z" fill="#D7D7D7" fill-opacity="0.38"/>  
  </g>
</svg>

<style>
  svg {
    transition: all ease 1s;
    cursor: pointer;
    animation: show 2s linear forwards;
  }

  @keyframes show {
    0% {
      opacity: 0%;
      
    }
    100% {
      opacity: 100%;
    }
  }

  svg:has(#granit:hover) {
    filter: drop-shadow(0 0 20px #147924);
    rotate: 2deg;
  }

  svg:has(#proton:hover) {
    filter: drop-shadow(0 0 20px #8311a0);
    rotate: -1deg;
  }

  svg:has(#wave:hover) {
    filter: drop-shadow(0 0 20px #0D65A4);
    rotate: 2deg;
  }

  svg:has(#rune:hover) {
    filter: drop-shadow(0 0 20px #9e7b1c);
    rotate: -1deg;
  }
</style>