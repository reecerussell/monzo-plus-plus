import React, { useState, useEffect } from "react";
import { Fetch } from "../../utils/fetch";
import Edit from "../../components/roles/edit";

const defaultData = {
	name: "",
};

const EditContainer = ({ id }) => {
	const [isMounted, setIsMounted] = useState(false);
	const [formData, setFormData] = useState(defaultData);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState(null);
	const [success, setSuccess] = useState(false);

	const handleFetch = async () => {
		if (loading) {
			return;
		}

		setLoading(true);

		await Fetch(
			`api/auth/roles/${id}`,
			null,
			async (res) => setFormData(await res.json()),
			setError
		);

		setLoading(false);
	};

	const handleSubmit = async (e) => {
		e.preventDefault();

		if (loading) {
			return;
		}

		setLoading(true);

		await Fetch(
			"api/auth/roles",
			{
				method: "PUT",
				body: JSON.stringify(formData),
			},
			() => setSuccess(true),
			setError
		);

		setLoading(false);
	};

	const handleFormUpdate = (e) => {
		const { name, value } = e.target;
		const data = { ...formData };
		data[name] = value;
		setFormData(data);
	};

	useEffect(() => {
		if (!isMounted) {
			handleFetch();
			setIsMounted(true);
		}

		return () => {};
	}, [isMounted, handleFetch]);

	useEffect(() => {
		if (success) {
			setTimeout(() => setSuccess(false), 4000);
		}
	}, [success]);

	return (
		<Edit
			id={id}
			data={formData}
			error={error}
			success={success}
			loading={loading}
			handleUpdate={handleFormUpdate}
			handleSubmit={handleSubmit}
			updateError={setError}
		/>
	);
};

export default EditContainer;
