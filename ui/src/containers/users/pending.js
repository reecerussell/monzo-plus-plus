import React, { useState, useEffect } from "react";
import Fetch from "../../utils/fetch";
import Pending from "../../components/users/pending";

const PendingContainer = () => {
	const [users, setUsers] = useState([]);
	const [error, setError] = useState(null);
	const [loading, setLoading] = useState(false);

	const handleFetchUsers = async () => {
		const res = await Fetch("http://localhost:9789/auth/users/pending");

		if (res.status === 200) {
			const data = await res.json();
			setUsers(data);
			setError(null);
		} else {
			if (res.headers.get("Content-Type") === "application/json") {
				const data = await res.json();
				setError(data.error);
			} else {
				setError(res.statusText);
			}
		}
	};

	const handleDelete = async (id, e) => {
		if (e) {
			e.preventDefault();
		}

		if (loading) {
			return;
		}

		setLoading(true);

		const res = await Fetch("http://localhost:9789/auth/users/" + id, {
			method: "DELETE",
		});

		if (res.status === 200) {
			await handleFetchUsers();
		} else {
			if (res.headers.get("Content-Type") === "application/json") {
				const data = await res.json();
				setError(data.error);
			} else {
				setError(res.statusText);
			}
		}

		setLoading(false);
	};

	const handleEnable = async (id, e) => {
		if (e) {
			e.preventDefault();
		}

		if (loading) {
			return;
		}

		setLoading(true);

		const res = await Fetch(
			"http://localhost:9789/auth/users/enable/" + id,
			{
				method: "POST",
			}
		);

		if (res.status === 200) {
			await handleFetchUsers();
		} else {
			if (res.headers.get("Content-Type") === "application/json") {
				const data = await res.json();
				setError(data.error);
			} else {
				setError(res.statusText);
			}
		}

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
