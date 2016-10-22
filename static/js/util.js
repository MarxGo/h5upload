function getFileSlice(file,start,length) {
    if(file.slice){
        return file.slice(start,length);
    } else if(file.mozSlice){
        // for FF
        return file.mozSlice(start,length);
    } else if(file.webkitSlice){
        // for Chrome
        return file.webkitSlice(start,length);
    } else {
     return null;
    }
}

function output(info) {
    document.getElementById("output").value += "\n" + info;
}
