import React from "react";
import { Header, Loader, Tab, Message, List } from "semantic-ui-react";

const Roles = ({
	roles,
	availableRoles,
	handleAddRole,
	handleRemoveRole,
	loading,
	error,
}) => {
	const assignedTab = () => (
		<List divided relaxed>
			{roles.map((p, key) => (
				<List.Item key={key}>
					<List.Icon
						name="group"
						size="large"
						verticalAlign="middle"
					/>
					<List.Content>
						<List.Header as="span">{p.name}</List.Header>
						<List.Description>
							<a onClick={() => handleRemoveRole(p.id)}>
								Click to remove role
							</a>
						</List.Description>
					</List.Content>
				</List.Item>
			))}
		</List>
	);

	const availableTab = () => (
		<List divided relaxed>
			{availableRoles.map((p, key) => (
				<List.Item key={key}>
					<List.Icon
						name="group"
						size="large"
						verticalAlign="middle"
					/>
					<List.Content>
						<List.Header as="span">{p.name}</List.Header>
						<List.Description>
							<a onClick={() => handleAddRole(p.id)}>
								Click to add role
							</a>
						</List.Description>
					</List.Content>
				</List.Item>
			))}
		</List>
	);

	return (
		<>
			<Loader active={loading} />

			<Header as="h2">Roles</Header>

			{error ? (
				<Message error header="An error occured!" content={error} />
			) : null}

			<Tab
				menu={{ secondary: true, pointing: true }}
				panes={[
					{ menuItem: "Assigned", render: assignedTab },
					{ menuItem: "Available", render: availableTab },
				]}
			/>
		</>
	);
};

export default Roles;
