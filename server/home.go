package server

const homehtml string = `
<!DOCTYPE html>
<html>
<head>
<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, minimum-scale=1, user-scalable=no, minimal-ui">
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


table {
  
  border-collapse: collapse;
  width: 100%;
  margin-bottom: 30px;
}

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

td, th {
  border: 1px solid #dddddd;
  text-align: left;
  padding: 8px;
}

tr:nth-child(even) {
  background-color: #dddddd;
}

.ten {
  width: 10%;
}


</style>
<title >Admin</title>
</head>
  <body >

    <div class="topnav">
      <a class="active" href="#home">Home</a>
      <a href="/about">About</a>
    </div> 
  
  <div id="container">
    
    <h1 style="text-align:center;">Monitored Inventory</h1>

  <p id="noneFoundMessage" style="display: none;color:red">No Items Found. Add some url's to settings.json</p>
  <table id="tab">
    <tr>
      
      <th>Name</th>
      <th>In Stock</th>
    </tr>
  
     {{range .Data}}
     <tr>
       <td>
         <a href="{{.Url}}"> {{.Name}} <a></a>
       </td>
       <td>
        {{.InStock}}
       </td>
      
    </tr>
    {{end}}
  </table>

  </div>
 
  
  <div id="footer">Caleb McCarthy 2021.</div>
</body>

</html>


<script>
function SetTableDisplay(bool){
  if (bool) {
    table = document.getElementById("tab");
    table.style.display = "table"
    document.getElementById("noneFoundMessage").style.display = "none";
  }else{
    table = document.getElementById("tab");
    table.style.display = "none"
    document.getElementById("noneFoundMessage").style.display = "inline";
  }
}

let tableData = {data:JSON.parse('{{.DataJson}}')}


function ClearTable() {
    var myTable = document.getElementById('tab');
    var rowCount = myTable.rows.length; 
    while(--rowCount) myTable.deleteRow(rowCount);
}

function CreateTable(inStock) {
    if(inStock.data.length <1){
      SetTableDisplay(false);
    }else{
      SetTableDisplay(true);
    }
    var table = document.getElementById("tab");

    inStockCount = 0;

    for(var i = 0; i < inStock.data.length; ++i) {
        var row = table.insertRow(table.rows.length);
        var cell1 = row.insertCell(0);
        var cell2 = row.insertCell(1);
        

        cell1.innerHTML = '<a href="'+inStock.data[i].url+'">'+ inStock.data[i].name+'<a>';
        cell2.innerHTML = inStock.data[i].instock
        if(inStock.data[i].instock ){
          ++inStockCount;
          alert(inStock.data[i].name+" in stock!\n"+inStock.data[i].url);
        }
        
    }
    if(inStockCount > 0) {
      document.title = "("+inStockCount+")"+" Admin Page"; 
    }else{
      document.title = "Admin Page"; 
    }
}

function compareTwo(oldTable,newTable){
  if(!oldTable.data || !newTable.data|| oldTable.data.length !== newTable.data.length){
    return true;
  }
  oldTable.data.sort();
  newTable.data.sort();
  for (i = 0;i < oldTable.data.length; ++i){
    if(oldTable.data[i].instock != newTable.data[i].instock){
      return true
    }
  }
  return false
}

function UpdateTable(){

fetch(window.location.href+'/api/items')
  .then(response => response.json())
  .then(data => {
    compareTwo
    
    // only Update Table if it has changed
    if(compareTwo(tableData,data)){

      tableData = data
      
      ClearTable();
      CreateTable(data);
    }

    
  });

  setTimeout(UpdateTable, 2000);
}
  
UpdateTable();

</script>

`
