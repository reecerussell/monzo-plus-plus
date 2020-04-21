import PropTypes from "prop-types";
import React, { Component } from "react";
import {
	Button,
	Container,
	Divider,
	Grid,
	Header,
	Icon,
	Image,
	List,
	Menu,
	Responsive,
	Segment,
	Sidebar,
	Visibility,
} from "semantic-ui-react";
import { Link } from "react-router-dom";

const getWidth = () => {
	const isSSR = typeof window === "undefined";

	return isSSR ? Responsive.onlyTablet.minWidth : window.innerWidth;
};

const HomepageHeading = ({ mobile }) => (
	<Container text>
		<Header
			as="h1"
			content="Monzo++"
			inverted
			style={{
				fontSize: mobile ? "2em" : "4em",
				fontWeight: "normal",
				marginBottom: 0,
				marginTop: mobile ? "1.5em" : "3em",
			}}
		/>
		<Header
			as="h2"
			content="Adding enhanced features to Monzo"
			inverted
			style={{
				fontSize: mobile ? "1.5em" : "1.7em",
				fontWeight: "normal",
				marginTop: mobile ? "0.5em" : "1.5em",
			}}
		/>
		<Button primary size="huge" as={Link} to="/register">
			Get Started
			<Icon name="right arrow" />
		</Button>
	</Container>
);

HomepageHeading.propTypes = {
	mobile: PropTypes.bool,
};

class DesktopContainer extends Component {
	state = {};

	hideFixedMenu = () => this.setState({ fixed: false });
	showFixedMenu = () => this.setState({ fixed: true });

	render() {
		const { children } = this.props;
		const { fixed } = this.state;

		return (
			<Responsive
				getWidth={getWidth}
				minWidth={Responsive.onlyTablet.minWidth}
			>
				<Visibility
					once={false}
					onBottomPassed={this.showFixedMenu}
					onBottomPassedReverse={this.hideFixedMenu}
				>
					<Segment
						inverted
						textAlign="center"
						style={{ minHeight: 700, padding: "1em 0em" }}
						vertical
					>
						<Menu
							fixed={fixed ? "top" : null}
							inverted={!fixed}
							pointing={!fixed}
							secondary={!fixed}
							size="large"
						>
							<Container>
								<Menu.Item as="a" active>
									Home
								</Menu.Item>
								<Menu.Item as="a" href="/docs">
									Docs
								</Menu.Item>
								<Menu.Item position="right">
									<Button
										as={Link}
										to="/login"
										inverted={!fixed}
									>
										Log in
									</Button>
									<Button
										as="a"
										href="https://github.com/reecerussell/monzo-plus-plus"
										inverted={!fixed}
										primary={fixed}
										style={{ marginLeft: "0.5em" }}
									>
										<Icon name="github" />
										Source
									</Button>
								</Menu.Item>
							</Container>
						</Menu>
						<HomepageHeading />
					</Segment>
				</Visibility>

				{children}
			</Responsive>
		);
	}
}

DesktopContainer.propTypes = {
	children: PropTypes.node,
};

class MobileContainer extends Component {
	state = {};

	handleSidebarHide = () => this.setState({ sidebarOpened: false });

	handleToggle = () => this.setState({ sidebarOpened: true });

	render() {
		const { children } = this.props;
		const { sidebarOpened } = this.state;

		return (
			<Responsive
				as={Sidebar.Pushable}
				getWidth={getWidth}
				maxWidth={Responsive.onlyMobile.maxWidth}
			>
				<Sidebar
					as={Menu}
					animation="push"
					inverted
					onHide={this.handleSidebarHide}
					vertical
					visible={sidebarOpened}
				>
					<Menu.Item as="a" active>
						Home
					</Menu.Item>
					<Menu.Item as="a" href="/docs">
						Docs
					</Menu.Item>
					<Menu.Item
						as="a"
						href="https://github.com/reecerussell/monzo-plus-plus"
					>
						View Source
					</Menu.Item>
					<Menu.Item>
						<Link to="/register">Log In</Link>
					</Menu.Item>
				</Sidebar>

				<Sidebar.Pusher dimmed={sidebarOpened}>
					<Segment
						inverted
						textAlign="center"
						style={{ minHeight: 350, padding: "1em 0em" }}
						vertical
					>
						<Container>
							<Menu inverted pointing secondary size="large">
								<Menu.Item onClick={this.handleToggle}>
									<Icon name="sidebar" />
								</Menu.Item>
								<Menu.Item position="right">
									<Button as={Link} to="/login" inverted>
										Log In
									</Button>
									<Button
										as="a"
										href="https://github.com/reecerussell/monzo-plus-plus"
										inverted
										style={{ marginLeft: "0.5em" }}
									>
										<Icon name="github" />
										Source
									</Button>
								</Menu.Item>
							</Menu>
						</Container>
						<HomepageHeading mobile />
					</Segment>

					{children}
				</Sidebar.Pusher>
			</Responsive>
		);
	}
}

MobileContainer.propTypes = {
	children: PropTypes.node,
};

const ResponsiveContainer = ({ children }) => (
	<div>
		<DesktopContainer>{children}</DesktopContainer>
		<MobileContainer>{children}</MobileContainer>
	</div>
);

ResponsiveContainer.propTypes = {
	children: PropTypes.node,
};

const HomepageLayout = () => (
	<ResponsiveContainer>
		<Segment style={{ padding: "8em 0em" }} vertical>
			<Grid container stackable verticalAlign="middle">
				<Grid.Row>
					<Grid.Column width={8}>
						<Header as="h3" style={{ fontSize: "2em" }}>
							Extendable
						</Header>
						<p style={{ fontSize: "1.33em" }}>
							Using plugins, users can add lots of different
							features to their account. Custom plugins can be
							built to extend the abilties of the Monzo++ system.
						</p>
						<Header as="h3" style={{ fontSize: "2em" }}>
							Open Source
						</Header>
						<p style={{ fontSize: "1.33em" }}>
							Fully open-source and licensed under the MIT
							license. Fork the Github repository to add your own
							plugins!
						</p>
					</Grid.Column>
					<Grid.Column floated="right" width={6}>
						<Image
							bordered
							rounded
							size="large"
							src="https://miro.medium.com/max/2000/1*8bPiDNL1K1ZdK9O_T5IVKw.png"
						/>
					</Grid.Column>
				</Grid.Row>
				<Grid.Row>
					<Grid.Column textAlign="center">
						<Button size="huge">
							<Icon name="github" /> View Source
						</Button>
					</Grid.Column>
				</Grid.Row>
			</Grid>
		</Segment>

		<Segment style={{ padding: "0em" }} vertical>
			<Grid celled="internally" columns="equal" stackable>
				<Grid.Row textAlign="center">
					<Grid.Column
						style={{ paddingBottom: "5em", paddingTop: "5em" }}
					>
						<Header as="h3" style={{ fontSize: "2em" }}>
							Custom built microservices
						</Header>
						<p style={{ fontSize: "1.33em" }}>
							Entirely built on containerised, custom-built
							microservices, written in Go.
						</p>
					</Grid.Column>
					<Grid.Column
						style={{ paddingBottom: "5em", paddingTop: "5em" }}
					>
						<Header as="h3" style={{ fontSize: "2em" }}>
							Scalable services
						</Header>
						<p style={{ fontSize: "1.33em" }}>
							All services are built to be resilient and scalable.
						</p>
					</Grid.Column>
				</Grid.Row>
			</Grid>
		</Segment>

		<Segment style={{ padding: "8em 0em" }} vertical>
			<Container text>
				<Header as="h3" style={{ fontSize: "2em" }}>
					Webhooks
				</Header>
				<p style={{ fontSize: "1.33em" }}>
					Using Monzo's webhook, plugins will listen for events on a
					Monzo account, then act accordingly. Currently, everytime a
					transaction is created, Monzo will trigger a webhook.
				</p>

				<Divider
					as="h4"
					className="header"
					horizontal
					style={{ margin: "3em 0em", textTransform: "uppercase" }}
				>
					More?
				</Divider>

				<Header as="h3" style={{ fontSize: "2em" }}>
					Security
				</Header>
				<p style={{ fontSize: "1.33em" }}>
					A user's Monzo API credentials are never exposed to the
					client, and each user's Monzo++ account is secured using
					short-lived OAuth2 tokens that are signed with RSA encoded
					asymmetric keys.
				</p>
			</Container>
		</Segment>

		<Segment inverted vertical style={{ padding: "5em 0em" }}>
			<Container textAlign="center">
				<Grid divided inverted stackable>
					<Grid.Column width={3}>
						<Header inverted as="h4" content="Info" />
						<List link inverted>
							<List.Item as="a" href="/">
								Home
							</List.Item>
						</List>
					</Grid.Column>
					<Grid.Column width={3}>
						<Header inverted as="h4" content="Developer" />
						<List link inverted>
							<List.Item as="a" href="/docs">
								API Docs
							</List.Item>
							<List.Item
								as="a"
								href="https://github.com/reecerussell/monzo-plus-plus"
							>
								Source
							</List.Item>
							<List.Item
								as="a"
								href="https://reece-russell.co.uk"
							>
								Reece Russell
							</List.Item>
						</List>
					</Grid.Column>
					<Grid.Column width={10}>
						<Header inverted as="h4" content="About" />
						<p>
							Monzo++ is an experimental project, built to explore
							and make use of webhooks through a custom third
							party system. Although experimental, Monzo++ is
							completely functional and can provide users with
							additional features.
						</p>
					</Grid.Column>
				</Grid>

				<Divider inverted section />
				<Image
					centered
					size="mini"
					src="https://lh3.googleusercontent.com/iDeb12CKMVdgDqBD9yJ9UehaWkKXFdPMtuUA8Jt0sOvxXzOm21qNGbA6D5_gdDZtAk4=w300"
				/>
				<List horizontal inverted divided link size="small">
					<List.Item as="a" href="http://reece-russell.co.uk">
						Designed and developed by{" "}
						<span style={{ textDecoration: "underline" }}>
							Reece Russell
						</span>
					</List.Item>
				</List>
			</Container>
		</Segment>
	</ResponsiveContainer>
);

export default HomepageLayout;
