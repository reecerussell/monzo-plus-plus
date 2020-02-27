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

		try {
			const res = await Fetch(
				"http://localhost:9789/auth/users?term=" + searchTerm
			);

			if (res.ok) {
				const data = await res.json();

				if (res.status == 200) {
					setUsers(data);
					setError(null);
				} else {
					setError(data.error);
				}
			} else {
				setError(res.statusText);
			}
		} catch {
			setError(
				"It seems like you don't have connection to the internet. Try again later!"
			);
		}

		setLoading(false);
	};

	const search = async e => {
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
