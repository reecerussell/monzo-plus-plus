import React, { useEffect, useState } from "react";
import { Fetch } from "../../../utils/fetch";
import Preferences from "../../../components/plugins/budget/preferences";

const PreferencesContainer = ({ id }) => {
	const [monthlyBudget, setMonthlyBudget] = useState("");
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState(null);
	const [success, setSuccess] = useState(null);
	const [isModalOpen, setModalOpen] = useState(false);

	const fetchData = async () => {
		if (loading) {
			return;
		}

		setLoading(true);

		await Fetch(
			"api/plugin/budget/preferences?userId=" + id,
			null,
			async (res) => {
				const { monthlyBudget } = await res.json();
				setMonthlyBudget(monthlyBudget / 100);
			},
			setError
		);

		setLoading(false);
	};

	const handleUpdate = async (e) => {
		if (e) {
			e.preventDefault();
		}

		if (loading) {
			return;
		}

		setLoading(true);

		await Fetch(
			"api/plugin/budget/preferences",
			{
				method: "PUT",
				body: JSON.stringify({
					userId: id,
					monthlyBudget: Math.ceil(parseFloat(monthlyBudget) * 100),
				}),
			},
			async () => {
				setError(null);
				setSuccess("Budget preferences saved!");
			},
			setError
		);

		setLoading(false);
	};

	const handleUpdateBudget = (e) => setMonthlyBudget(e.target.value);

	const toggleModal = () => setModalOpen(!isModalOpen);

	useEffect(() => {
		fetchData();
	}, []);

	useEffect(() => {
		if (success !== null) {
			setTimeout(() => setSuccess(null), 3000);
		}

		if (error !== null) {
			setTimeout(() => setSuccess(null), 5000);
		}
	}, [success, error]);

	return (
		<Preferences
			error={error}
			success={success}
			loading={loading}
			handleUpdateBudget={handleUpdateBudget}
			handleUpdate={handleUpdate}
			monthlyBudget={monthlyBudget}
			showModal={isModalOpen}
			toggleModal={toggleModal}
		/>
	);
};

export default PreferencesContainer;
