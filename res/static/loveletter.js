function playcard(position, card) {
	document.getElementById('cards').value = position;
	document.getElementById('targets').value = 'computer';
	switch (card.toLowerCase()) {
		case 'guard':
			$('#guardmodal').modal();
			break;
		case 'prince':
			$('#princemodal').modal();
			break;
		default:
			submit();
	}
}

function playguard(card) {
	document.getElementById('guess').value = card;
	submit()
}

function playprince(target) {
	document.getElementById('targets').value = target;
	submit()
}

function submit() {
	document.getElementById('playform').submit();
}
