window.onload = function(){
    var form = document.getElementById("the-form");
    var city = document.getElementById("city");

    form.addEventListener("submit", function(event){
        event.preventDefault();

        fetch("http://localhost/zips/city/" + city.value)
            .then(function(response){
                console.log(response);
                return response.json();
            })
            .then(function(data){
                var output = document.createElement("p");
                console.log(data);
                output.textContent = JSON.stringify(data);
                document.body.appendChild(output);
            })
            .catch(function(err){
                console.error(err);
            })

    })
}
