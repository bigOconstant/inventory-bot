package server

const abouthtml = `
<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, minimum-scale=1, user-scalable=no, minimal-ui">

<title>About</title>

<style>
    /*

menu

*/
/* Add a black background color to the top navigation */
.topnav {
  background-color: rgb(51, 51, 51);
  overflow: hidden;
  margin-left: -8px;
  margin-right: -8px;
  margin-top:-8px;
  
}

/* Style the links inside the navigation bar */
.topnav a {
  float: left;
  color: #f2f2f2;
  text-align: center;
  padding: 14px 16px;
  text-decoration: none;
  font-size: 17px;
}

/* Change the color of links on hover */
.topnav a:hover {
  background-color: rgb(232, 225, 247);
  color: black;
}

/* Add a color to the active/current link */
.topnav a.active {
  background-color: hsl(253, 39%, 49%);
  color: white;
}

/* end menu*/

html, body {
    
    font-family: Arial, Helvetica, sans-serif;
    height: 100%;
    text-align: center;
}

#container {
    
    margin-bottom: 2em;
    min-height: 100%;
    overflow: auto;
    padding: 0 1em;
    text-align: justify;
}

#footer {
    bottom: 0;
    color: #707070;
    height: 2em;
    left: 0;
    position: relative;
    font-size: small;
    width:100%;
}

a:link {
  text-decoration: none;
}

a:visited {
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}

a:active {
  text-decoration: underline;
}
</style>
</head>
<div class="topnav">
    <a  href="/">Home</a>
    <a class="active" >About</a>
</div> 
<body>
    
   
    <div id="container">
    <h1 style="text-align:center;">Inventory Bot</h1>

    <p style="text-align:center;">Scans a given url for a Add to Cart Phrase and reports back whether it is found.</p>

    <h2 style="text-align:center;">Configuration</h2>

    <p style="text-align:center;">Use the settings.json file to add or remove items, configure a port, or add a discord web hook to send notifications</p>

    <h3 style="text-align:center;">Github</h3>

    <a style="text-align:center; width:100%" href="https://github.com/camccar/GPU-bot"><p style="text-align:center; width:100%">https://github.com/camccar/GPU-bot</p></a>

    
</div>
    <div id="footer" style="text-align:center;" >Caleb McCarthy 2021.</div>
</body>
</html> 
`
