import React, { useState, useEffect } from "react";
import Fetch from "../../utils/fetch";
import Pending from "../../components/users/pending";

const PendingContainer = () => {
	const [users, setUsers] = useState([]);
	const [error, setError] = useState(null);
	const [loading, setLoading] = useState(false);

	const handleFetchUsers = async () => {
		if (loading) {
			return;
		}

		setLoading(true);

		await Fetch(
			"api/auth/users/pending",
			null,
			async (res) => {
				setError(null);
				setUsers(await res.json());
			},
			setError
		);

		setLoading(false);
	};

	const handleDelete = async (id, e) => {
		if (e) {
			e.preventDefault();
		}

		if (loading) {
			return;
		}

		setLoading(true);

		await Fetch(
			"api/auth/users/enable/" + id,
			{
				method: "DELETE",
			},
			handleFetchUsers,
			setError
		);

		setLoading(false);
	};

	const handleEnable = async (id, e) => {
		if (e) {
			e.preventDefault();
		}

		if (loading) {
			return;
		}

		await Fetch(
			"api/auth/users/enable/" + id,
			{
				method: "POST",
			},
			handleFetchUsers,
			setError
		);

		setLoading(false);
	};

	const fetchUsers = async () => {
		if (loading) {
			return;
		}

		setLoading(true);

		await handleFetchUsers();

		setLoading(false);
	};

	useEffect(() => {
		fetchUsers();
	}, []);

	return (
		<Pending
			users={users}
			error={error}
			loading={loading}
			handleDelete={handleDelete}
			handleEnable={handleEnable}
		/>
	);
};

export default PendingContainer;
