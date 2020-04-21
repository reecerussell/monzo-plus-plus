import React, { useState } from "react";
import * as User from "../../utils/user";
import Fetch from "../../utils/fetch";
import Delete from "../../components/account/delete";

const DeleteContainer = () => {
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState(null);
	const [redirect, setRedirect] = useState(null);

	const handleDelete = async (e) => {
		e.preventDefault();

		if (loading) {
			return;
		}

		setLoading(true);

		await Fetch(
			"api/auth/users/" + User.GetId(),
			{
				method: "DELETE",
			},
			async () => {
				setError(null);
				User.Logout();
				setRedirect("/");
			},
			setError
		);

		setLoading(false);
	};

	return (
		<Delete
			handleSubmit={handleDelete}
			error={error}
			loading={loading}
			redirect={redirect}
		/>
	);
};

export default DeleteContainer;
