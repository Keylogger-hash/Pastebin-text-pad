function main(){
    var button = document.querySelector("button");
    var textarea = document.querySelector("textarea");
    button.addEventListener('click',function(){
        textarea.select()
        document.execCommand("copy")
    })
}
main()