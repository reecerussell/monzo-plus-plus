import React, { useState, useEffect } from "react";
import { Fetch } from "../../utils/fetch";
import List from "../../components/roles/list";

const ListContainer = () => {
	const [roles, setRoles] = useState([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState(null);
	const [searchTerm, setSearchTerm] = useState("");

	const handleFetch = async () => {
		if (loading) {
			return;
		}

		setLoading(true);

		await Fetch(
			"api/auth/roles?term=" + searchTerm,
			null,
			async (res) => {
				setError(null);
				setRoles(await res.json());
			},
			setError
		);

		setLoading(false);
	};

	const handleFlush = async () =>
		await Fetch(
			"api/auth/permissions/flush",
			{ method: "POST" },
			null,
			setError
		);

	const handleSearchUpdate = (e) => setSearchTerm(e.target.value);

	useEffect(() => {
		handleFetch();
	}, []);

	return (
		<List
			roles={roles}
			loading={loading}
			error={error}
			searchTerm={searchTerm}
			updateSearchTerm={handleSearchUpdate}
			onSearch={handleFetch}
			handleFlush={handleFlush}
		/>
	);
};

export default ListContainer;
