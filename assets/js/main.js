let form = document.querySelector('.js-form'),
	formfio = document.querySelectorAll('.js-fio');
	
	form.onsubmit = function(){
		console.log('works');
		let emptyfio = Array.from(formfio).filter(inputMyWebApp => inputMyWebApp.value ==="");
		
		formfio.forEach(function(inputMyWebApp) {
			if (inputMyWebApp.value === ""){
				inputMyWebApp.classList.add('error');
				
				alert('Вы не заполнили тему');
			} else {
				inputMyWebApp.classList.remove('error');
			}
		});
		
		if(emptyfio.length !== 0){
			console.log('input not filled');
			return false;
		}
		
		
	}	
  
	document.getElementById("foo").onkeypress = function(e) {
		var chr = String.fromCharCode(e.which);
		if ("12345NOABC".indexOf(chr) < 0)
			return false;
	};