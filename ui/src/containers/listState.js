import { useState, useEffect } from "react";
import { Fetch, BaseUrl } from "../utils/fetch";

const ListState = (path) => {
	const [error, setError] = useState(null);
	const [loading, setLoading] = useState(false);
	const [searchTerm, setSearchTerm] = useState("");
	const [items, setItems] = useState([]);

	const handleFetch = async () => {
		if (loading) {
			return;
		}

		setLoading(true);

		const url = new URL(path, BaseUrl);
		url.searchParams.set("term", searchTerm);

		await Fetch(
			url.toString(),
			null,
			async (res) => setItems(await res.json()),
			setError,
			""
		);

		setLoading(false);
	};

	useEffect(() => {
		handleFetch();
	}, []);

	return [items, loading, error, searchTerm, setSearchTerm, handleFetch];
};

export default ListState;
