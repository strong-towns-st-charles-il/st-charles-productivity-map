const puppeteer = require('puppeteer');
const fs = require('fs');

async function scrapeParcel(url) {
  const browser = await puppeteer.launch();
  const page = await browser.newPage();
  await page.goto(url);

  const [el1] = await page.$x('//*[@id="parcel-views"]/div[3]/table/tbody/tr[5]/td[3]/div[3]');
  const txt1 = await el1.getProperty('textContent');
  const tax = await txt1.jsonValue();

  const [el2] = await page.$x('//*[@id="parcel-views"]/div[3]/table/tbody/tr[6]/td[2]/div[2]');
  const txt2 = await el2.getProperty('textContent');
  const acres = await txt2.jsonValue();

  console.log({tax, acres})
  browser.close();
}
//scrapeParcel('https://kaneil.devnetwedge.com/parcel/view/0927391001/2021')


// Data from: https://kaneil.devnetwedge.com/
baseURL = 'https://kaneil.devnetwedge.com/parcel/view/'
getData(baseURL);

async function getData(baseURL) {
  fs.readFile('data/sample-parcel.csv', 'utf-8', (err, data) => {
    if (err) console.log(err);
    else {
      const table = data.split(/\n/).slice(1,-1);
      table.forEach(row => {
        const column = CSVtoArray(row);
        const year = column[0];
        const parcel = removeHyphen(column[1]);
        url = baseURL+parcel+'/'+year;
        scrapeParcel(url);
      });
    }
  });
}

function removeHyphen(str) {
  const a = str.split('-')
  return a.join('');

}

// https://stackoverflow.com/questions/8493195/how-can-i-parse-a-csv-string-with-javascript-which-contains-comma-in-data
function CSVtoArray(text) {
  var re_valid = /^\s*(?:'[^'\\]*(?:\\[\S\s][^'\\]*)*'|"[^"\\]*(?:\\[\S\s][^"\\]*)*"|[^,'"\s\\]*(?:\s+[^,'"\s\\]+)*)\s*(?:,\s*(?:'[^'\\]*(?:\\[\S\s][^'\\]*)*'|"[^"\\]*(?:\\[\S\s][^"\\]*)*"|[^,'"\s\\]*(?:\s+[^,'"\s\\]+)*)\s*)*$/;
  var re_value = /(?!\s*$)\s*(?:'([^'\\]*(?:\\[\S\s][^'\\]*)*)'|"([^"\\]*(?:\\[\S\s][^"\\]*)*)"|([^,'"\s\\]*(?:\s+[^,'"\s\\]+)*))\s*(?:,|$)/g;
  // Return NULL if input string is not well formed CSV string.
  if (!re_valid.test(text)) return null;
  var a = [];                     // Initialize array to receive values.
  text.replace(re_value, // "Walk" the string using replace with callback.
    function(m0, m1, m2, m3) {
      // Remove backslash from \' in single quoted values.
      if      (m1 !== undefined) a.push(m1.replace(/\\'/g, "'"));
        // Remove backslash from \" in double quoted values.
      else if (m2 !== undefined) a.push(m2.replace(/\\"/g, '"'));
      else if (m3 !== undefined) a.push(m3);
      return ''; // Return empty string.
    });
  // Handle special case of empty last value.
  if (/,\s*$/.test(text)) a.push('');
  return a;
};
