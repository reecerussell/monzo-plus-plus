import React, { useState } from "react";
import { Fetch } from "../../utils/fetch";
import Create from "../../components/roles/create";

const defaultData = {
	name: "",
};

const CreateContainer = () => {
	const [formData, setFormData] = useState(defaultData);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState(null);
	const [redirect, setRedirect] = useState(null);

	const handleSubmit = async e => {
		e.preventDefault();

		if (loading) {
			return;
		}

		setLoading(true);

		await Fetch(
			"auth/roles",
			{
				method: "POST",
				body: JSON.stringify(formData),
			},
			async res => {
				const { id } = await res.json();
				setRedirect(`/roles/details/${id}`);
			},
			setError
		);

		setLoading(false);
	};

	const handleFormUpdate = e => {
		const { name, value } = e.target;
		const data = { ...formData };
		data[name] = value;
		setFormData(data);
	};

	return (
		<Create
			data={formData}
			error={error}
			loading={loading}
			redirect={redirect}
			handleUpdate={handleFormUpdate}
			handleSubmit={handleSubmit}
		/>
	);
};

export default CreateContainer;
