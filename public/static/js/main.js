var button = document.querySelector("#edit-button")
button.onclick = function(){
    var textarea = document.querySelector("textarea");
    document.getElementById("edit-hidden").type='submit';
    textarea.readOnly = false;
    alert("Now you can edit paste!")    
}
function main(){
    var button = document.querySelector("button");
    var textarea = document.querySelector("textarea");
    button.addEventListener('click',function(){
        textarea.select()
        document.execCommand("copy")
    })
}
main()