import React, { useState, useEffect } from "react";

function App() {
	const [query, setQuery] = useState("");
	const [results, setResults] = useState(null);

	useEffect(() => {
		if (query !== "") {
			const fetchData = async () => {
				const response = await fetch(
					`/api/autocomplete?query=${query}`,
				);
				const data = await response.json();
				const processedData = data.map((resultGroup) =>
					resultGroup.map((result) => ({
						name: result.name,
						role: result.role.join(", "),
					}))
				);
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
						{results && results.flat().map((result, index) => (
							<li key={index}>
								<strong>Name:</strong> {result.name} <br />
								<strong>Role:</strong> {result.role}
							</li>
						))}
					</ul>
			</div>
		</div>
	);
}

export default App;
