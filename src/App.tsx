import React, { useState, useEffect, useCallback } from "react";

function App() {
	const [query, setQuery] = useState("");
	const [results, setResults] = useState(null);
	const [loading, setLoading] = useState(false);

	const debounce = (func, delay) => {
		let timeoutId;
		return (...args) => {
			if (timeoutId) {
				clearTimeout(timeoutId);
			}
			timeoutId = setTimeout(() => {
				func(...args);
			}, delay);
		};
	};

	const debouncedFetchData = useCallback(
		debounce(async (query) => {
			setLoading(true);
			const response = await fetch(`/api/autocomplete?query=${query}`);
			const data = await response.json();

			const processedData = data[0].map((item) => {
				const role = item.role[0] === "" ? item.role[1] : item.role[0];
				return {
					name: item.name,
					role: role,
				};
			});

			setResults(processedData);
			setLoading(false);
		}, 300),
		[]
	);

	useEffect(() => {
		if (query !== "") {
			debouncedFetchData(query);
		}
	}, [query, debouncedFetchData]);

	return (
		<div>
			<div>
				<label>
					Query:
					<input
						type="text"
						value={query}
						onChange={(e) => setQuery(e.target.value)}
					/>
				</label>
			</div>
			<div>
				<h3>Results:</h3>
				{loading ? (
					<p>Loading...</p>
				) : (
					<ul>
						{results &&
							results.map((result, index) => (
								<li key={index}>
									{result.name} - {result.role}
								</li>
							))}
					</ul>
				)}
			</div>
		</div>
	);
}

export default App;
