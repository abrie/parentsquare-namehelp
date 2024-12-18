import React, { useState, useEffect } from "react";

function App() {
	const [query, setQuery] = useState("");
	const [results, setResults] = useState(null);

	useEffect(() => {
		if (query !== "") {
			const fetchData = async () => {
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
			};
			fetchData();
		}
	}, [query]);

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
				<ul>
					{results &&
						results.map((result, index) => (
							<li key={index}>
								{result.name} - {result.role}
							</li>
						))}
				</ul>
			</div>
		</div>
	);
}

export default App;
