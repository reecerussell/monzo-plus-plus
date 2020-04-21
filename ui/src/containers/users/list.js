import React, { useState, useEffect } from "react";
import Fetch from "../../utils/fetch";
import List from "../../components/users/list";

const ListContainer = () => {
	const [users, setUsers] = useState([]);
	const [error, setError] = useState();
	const [loading, setLoading] = useState(false);
	const [searchTerm, updateSearchTerm] = useState("");

	const fetchUsers = async () => {
		if (loading) {
			return;
		}

		setLoading(true);

		await Fetch(
			"api/auth/users?term=" + searchTerm,
			null,
			async (res) => {
				setUsers(await res.json());
				setError(null);
			},
			setError
		);

		setLoading(false);
	};

	const search = async (e) => {
		if (e) {
			e.preventDefault();
		}

		await fetchUsers();
	};

	useEffect(() => {
		fetchUsers();
	}, []);

	return (
		<List
			users={users}
			loading={loading}
			error={error}
			searchTerm={searchTerm}
			updateSearchTerm={updateSearchTerm}
			onSearch={search}
		/>
	);
};

export default ListContainer;
