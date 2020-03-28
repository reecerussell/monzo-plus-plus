import React from "react";
import { Header, Loader, Tab, Message, List } from "semantic-ui-react";

const Permissions = ({
	permissions,
	availablePermissions,
	handleAddPermission,
	handleRemovePermission,
	loading,
	error,
}) => {
	const assignedTab = () => (
		<List divided relaxed>
			{permissions.map((p, key) => (
				<List.Item key={key}>
					<List.Icon name="key" size="large" verticalAlign="middle" />
					<List.Content>
						<List.Header as="span">{p.name}</List.Header>
						<List.Description>
							<a onClick={() => handleRemovePermission(p.id)}>
								Click to remove permission
							</a>
						</List.Description>
					</List.Content>
				</List.Item>
			))}
		</List>
	);

	const availableTab = () => (
		<List divided relaxed>
			{availablePermissions.map((p, key) => (
				<List.Item key={key}>
					<List.Icon name="key" size="large" verticalAlign="middle" />
					<List.Content>
						<List.Header as="span">{p.name}</List.Header>
						<List.Description>
							<a onClick={() => handleAddPermission(p.id)}>
								Click to add permission
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

			<Header as="h2">Permissions</Header>

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

export default Permissions;
