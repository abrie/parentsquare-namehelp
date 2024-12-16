import React, { useState, useEffect } from "react";

function App() {
	const [schoolId, setSchoolId] = useState(732);
	const [limit, setLimit] = useState(25);
	const [chat, setChat] = useState(1);
	const [query, setQuery] = useState("");
	const [results, setResults] = useState(null);

	useEffect(() => {
		if (query !== "") {
			const fetchData = async () => {
				const response = await fetch(
					`/api/autocomplete?school_id=${schoolId}&limit=${limit}&chat=${chat}&query=${query}`,
				);
				const data = await response.json();
				setResults(data);
			};
			fetchData();
		}
	}, [query, schoolId, limit, chat]);

	return (
		<div>
			<div>
				<label>
					School ID:
					<input
						type="number"
						value={schoolId}
						onChange={(e) => setSchoolId(Number(e.target.value))}
					/>
				</label>
			</div>
			<div>
				<label>
					Limit:
					<input
						type="number"
						value={limit}
						onChange={(e) => setLimit(Number(e.target.value))}
					/>
				</label>
			</div>
			<div>
				<label>
					Chat:
					<input
						type="number"
						value={chat}
						onChange={(e) => setChat(Number(e.target.value))}
					/>
				</label>
			</div>
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
