(function ($) {
	'use strict';
	$(function () {
		var todoListItem = $('.todo-list');
		var todoListInput = $('.todo-list-input');
		$('.todo-list-add-btn').on('click', function (event) {
			event.preventDefault();

			var item = $(this).prevAll('.todo-list-input').val();
			if (item) {
				$.ajax({
					url: '/api/todos',
					type: 'POST',
					contentType: 'application/json',
					data: JSON.stringify({ name: item, completed: false }),
					success: addItem,
				});
				todoListInput.val('');
			}
		});
		var addItem = (item) => {
			todoListItem.append(
				`<li ${item.completed ? "class='completed'" : ''} id='${
					item._id
				}'><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' ${
					item.completed ? "checked='checked'" : ''
				}/>${
					item.name
				}<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>`
			);
		};

		$.get('/api/todos', (items) => {
			for (let i = 0; i < items.length; i++) {
				addItem(items[i]);
			}
		});

		todoListItem.on('change', '.checkbox', function () {
			const id = $(this).closest('li').attr('id');
			const $self = $(this);

			let complete = true;
			if ($(this).attr('checked')) {
				complete = false;
			}
			$.get(`/api/complete/${id}?complete=${complete}`, function () {
				if (complete) {
					$self.attr('checked', 'checked');
				} else {
					$self.removeAttr('checked');
				}
				$self.closest('li').toggleClass('completed');
			});
		});

		todoListItem.on('click', '.remove', function (e) {
			const id = $(this).closest('li').attr('id');
			const $self = $(this);
			$.ajax({
				url: `/api/todos/${id}`,
				type: 'DELETE',
				success: function (data) {
					if (data) {
						const str = e.target.parentElement.firstChild.innerText;
						$self.parent().remove();
						alert(`You will delete Todo: ${str}`);
					}
				},
			});
		});
	});
})(jQuery);
