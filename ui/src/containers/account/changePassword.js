import React, { useState } from "react";
import Fetch from "../../utils/fetch";
import * as User from "../../utils/user";
import ChangePassword from "../../components/account/changePassword";

const defaultFormData = {
	currentPassword: {
		value: "",
		error: null,
	},
	newPassword: {
		value: "",
		error: null,
	},
	confirmPassword: {
		value: "",
		error: null,
	},
};

const ChangePasswordContainer = () => {
	const [formData, setFormData] = useState(defaultFormData);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState(null);
	const [success, setSuccess] = useState(null);

	const handleSubmit = async (e) => {
		e.preventDefault();

		if (
			loading ||
			formData.currentPassword.error ||
			formData.newPassword.error ||
			formData.confirmPassword.error
		) {
			return;
		}

		setLoading(true);

		try {
			const res = await Fetch(
				"http://localhost:9789/api/auth/users/changepassword",
				{
					method: "POST",
					body: JSON.stringify({
						currentPassword: formData.currentPassword.value,
						newPassword: formData.newPassword.value,
					}),
				}
			);

			if (res.status == 200) {
				setError(null);
				setFormData(defaultFormData);
				setSuccess("Password changed successfully!");
			} else {
				const data = await res.json();
				setError(data.error);
				setSuccess(null);
			}
		} catch {
			setError(
				"It seems like you don't have connection to the internet. Try again later!"
			);
			setSuccess(null);
		}

		setLoading(false);
	};

	const handleUpdateCurrentPassword = (e) => {
		const { value } = e.target;

		const data = formData.currentPassword;
		data.value = value;

		if (value === "" || value.length < 1) {
			data.error = "This field is required.";
		} else {
			data.error = null;
		}

		const newFormData = { ...formData };
		newFormData.currentPassword = data;

		setFormData(newFormData);
	};

	const handleUpdateNewPassword = (e) => {
		const { value } = e.target;

		const data = formData.newPassword;
		data.value = value;

		if (value === "" || value.length < 1) {
			data.error = "This field is required.";
		} else {
			data.error = null;
		}

		const newFormData = { ...formData };
		newFormData.newPassword = data;

		setFormData(newFormData);
	};

	const handleUpdateConfirmPassword = (e) => {
		const { value } = e.target;

		const data = formData.newPassword;
		data.value = value;

		if (value === "" || value.length < 1) {
			data.error = "This field is required.";
		} else if (value !== formData.newPassword.value) {
			data.error = "New passwords do not match.";
		} else {
			data.error = null;
		}

		const newFormData = { ...formData };
		newFormData.newPassword = data;

		setFormData(newFormData);
	};

	return (
		<ChangePassword
			data={formData}
			loading={loading}
			error={error}
			success={success}
			username={User.GetUsername()}
			handleSubmit={handleSubmit}
			handleUpdateCurrentPassword={handleUpdateCurrentPassword}
			handleUpdateNewPassword={handleUpdateNewPassword}
			handleUpdateConfirmPassword={handleUpdateConfirmPassword}
		/>
	);
};

export default ChangePasswordContainer;
