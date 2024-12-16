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
		<div className="p-4">
			<div className="mb-4">
				<label className="block text-sm font-medium text-gray-700">
					School ID:
					<input
						type="number"
						value={schoolId}
						onChange={(e) => setSchoolId(Number(e.target.value))}
						className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
					/>
				</label>
			</div>
			<div className="mb-4">
				<label className="block text-sm font-medium text-gray-700">
					Limit:
					<input
						type="number"
						value={limit}
						onChange={(e) => setLimit(Number(e.target.value))}
						className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
					/>
				</label>
			</div>
			<div className="mb-4">
				<label className="block text-sm font-medium text-gray-700">
					Chat:
					<input
						type="number"
						value={chat}
						onChange={(e) => setChat(Number(e.target.value))}
						className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
					/>
				</label>
			</div>
			<div className="mb-4">
				<label className="block text-sm font-medium text-gray-700">
					Query:
					<input
						type="text"
						value={query}
						onChange={(e) => setQuery(e.target.value)}
						className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
					/>
				</label>
			</div>
			<div>
				<h3 className="text-lg font-medium text-gray-900">Results:</h3>
				<pre className="mt-2 p-4 bg-gray-100 rounded-md">{results && JSON.stringify(results, null, 2)}</pre>
			</div>
		</div>
	);
}

export default App;
