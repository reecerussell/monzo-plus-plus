import React, { useEffect, useState } from "react";
import { Fetch } from "../../utils/fetch";
import PropTypes from "prop-types";
import Permissions from "../../components/roles/permissions";

const PermissionsContainer = ({ id }) => {
	const [perms, setPerms] = useState([]);
	const [available, setAvailable] = useState([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState(null);

	const handleFetchPerms = async () => {
		if (loading) {
			return;
		}

		setLoading(true);

		const permissions = Fetch(
			`api/auth/roles/permissions/${id}`,
			null,
			async (res) => setPerms(await res.json()),
			setError
		);

		const availablePermissions = Fetch(
			`api/auth/roles/availablePermissions/${id}`,
			null,
			async (res) => setAvailable(await res.json()),
			setError
		);

		await Promise.all([permissions, availablePermissions]);

		setLoading(false);
	};

	const handleAddPermission = async (permissionId) =>
		await Fetch(
			"api/auth/roles/permission",
			{
				method: "POST",
				body: JSON.stringify({
					permissionId: parseInt(permissionId),
					roleId: id,
				}),
			},
			handleFetchPerms,
			setError
		);

	const handleRemovePermission = async (permissionId) =>
		await Fetch(
			"api/auth/roles/permission",
			{
				method: "DELETE",
				body: JSON.stringify({
					permissionId: parseInt(permissionId),
					roleId: id,
				}),
			},
			handleFetchPerms,
			setError
		);

	useEffect(() => {
		handleFetchPerms();
	}, []);

	return (
		<Permissions
			handleAddPermission={handleAddPermission}
			handleRemovePermission={handleRemovePermission}
			permissions={perms}
			availablePermissions={available}
			error={error}
			loading={loading}
		/>
	);
};

export default PermissionsContainer;
