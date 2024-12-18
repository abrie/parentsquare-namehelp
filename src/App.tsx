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
				setResults(data);
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
				<pre>{results && JSON.stringify(results, null, 2)}</pre>
			</div>
		</div>
	);
}

export default App;
