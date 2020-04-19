import React, { useState, useEffect } from "react";
import { Fetch } from "../../utils/fetch";
import Roles from "../../components/users/roles";

const RolesContainer = ({ id }) => {
	const [roles, setRoles] = useState([]);
	const [available, setAvailable] = useState([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState(null);

	const handleFetch = async () => {
		if (loading) {
			return;
		}

		setLoading(true);

		const roles = Fetch(
			`api/auth/users/roles/${id}`,
			null,
			async (res) => setRoles(await res.json()),
			setError
		);

		const availableRoles = Fetch(
			`api/auth/users/availableRoles/${id}`,
			null,
			async (res) => setAvailable(await res.json()),
			setError
		);

		await Promise.all([roles, availableRoles]);

		setLoading(false);
	};

	const handleAddRole = async (roleId) =>
		await Fetch(
			"api/auth/users/roles",
			{
				method: "POST",
				body: JSON.stringify({
					roleId: roleId,
					userId: id,
				}),
			},
			handleFetch,
			setError
		);

	const handleRemoveRole = async (roleId) =>
		await Fetch(
			"api/auth/users/roles",
			{
				method: "DELETE",
				body: JSON.stringify({
					roleId: roleId,
					userId: id,
				}),
			},
			handleFetch,
			setError
		);

	useEffect(() => {
		handleFetch();
	}, []);

	return (
		<Roles
			roles={roles}
			availableRoles={available}
			handleAddRole={handleAddRole}
			handleRemoveRole={handleRemoveRole}
			error={error}
			loading={loading}
		/>
	);
};

export default RolesContainer;
