document.addEventListener("DOMContentLoaded", function () {
	const combobox = document.getElementById("combobox");
	const optionsList = document.getElementById("options");
	const items = optionsList.querySelectorAll("li");

	let selectedIndex = -1;

	// Toggle dropdown visibility
	combobox.addEventListener("click", function () {
		if (optionsList.style.display === "block") {
			optionsList.style.display = "none";
		} else {
			optionsList.style.display = "block";
		}
	});

	// Filter list when typing
	combobox.addEventListener("input", function () {
		const filter = combobox.value.toLowerCase();
		items.forEach((item) => {
			if (item.innerText.toLowerCase().includes(filter)) {
				item.style.display = "block";
			} else {
				item.style.display = "none";
			}
		});
	});

	// Handle keyboard navigation
	combobox.addEventListener("keydown", function (e) {
		switch (e.key) {
			case "ArrowDown":
				selectedIndex++;
				if (selectedIndex > items.length - 1) selectedIndex = 0;
				break;
			case "ArrowUp":
				selectedIndex--;
				if (selectedIndex < 0) selectedIndex = items.length - 1;
				break;
			case "Enter":
				if (selectedIndex >= 0) {
					selectItem(items[selectedIndex]);
				}
				return;
			case "Escape":
				optionsList.style.display = "none";
				return;
			default:
				return;
		}

		items.forEach((item, index) => {
			if (index === selectedIndex) {
				item.classList.add("text-white", "bg-indigo-600");
			} else {
				item.classList.remove("text-white", "bg-indigo-600");
				item.classList.add("text-gray-900");
			}
		});
	});

	items.forEach((item) => {
		// Highlight on hover
		item.addEventListener("mouseenter", function () {
			item.classList.add("text-white", "bg-indigo-600");
		});
		item.addEventListener("mouseleave", function () {
			item.classList.remove("text-white", "bg-indigo-600");
			item.classList.add("text-gray-900");
		});

		// Select on click
		item.addEventListener("click", function () {
			selectItem(item);
		});
	});

	function selectItem(item) {
		const itemId = item.id.split("-")[1];
		combobox.value = itemId;
		optionsList.style.display = "none";

		fetchCoins(itemId);
	}

	function fetchCoins(coinId, currency = "cad") {
		const URL = `/coins?currency=${currency}&ids=${coinId}`;

		fetch(URL)
			.then((response) => {
				// status code is 200-299
				if (response.ok) {
					window.location.href = URL;
				} else {
					return response.text().then((text) => {
						alert(`Error: ${text}`);
					});
				}
			})
			.catch((error) => {
				alert(`Network error: ${error.message}`);
			});
	}
});
