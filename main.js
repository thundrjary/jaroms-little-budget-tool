

function addRow() {
  var tableRef = document.getElementById('myTable');

  // Insert a row in the table at row index 0
  var newRow   = tableRef.insertRow(tableRef.rows.length);

  // Insert a cell in the row at index 0
  var c1  = newRow.insertCell(0);
  var c2  = newRow.insertCell(0);
  var c3  = newRow.insertCell(0);

  // Append a text node to the cell
  var item  = document.createTextNode('<span style="color: lightblue;">Item</span>')
  var dollar = document.createTextNode('$')
  var amount  = document.createTextNode('<span  style="color: lightblue;">Amount</span>')
  c1.appendChild(item);
  c2.appendChild(dollar);
  c3.appendChild(amount);
}
