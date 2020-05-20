pragma solidity 0.5.15;

contract ENS {

    // Domain holds info needed to keep track of ownership`
    struct Domain {
        address pointTo;
        address owner;
    }

    // domains holds all name: domain
    mapping(bytes32 => Domain) public domains;

    // ChangeDomain logs the alteration of a domain by its owner
    event ChangeDomain(bytes32 indexed name, address indexed pointer, address indexed owner);

    // AddDomains logs new domains
    event AddDomain(bytes32 indexed name, address indexed pointer, address indexed owner);

    // change alters an already existing domain pointer
    function change(bytes32 name, address addr) public {
        require(msg.sender == domains[name].owner, "user does not have rights to desired domain");
        domains[name].pointTo = addr;
        emit ChangeDomain(name, addr, msg.sender);
    }

    // add creates a new domain and assigns it to the provided address
    function add(bytes32 name, address addr) public {
        require(domains[name].pointTo == address(0x0), "domain already exists");
        domains[name].pointTo = addr;
        domains[name].owner = msg.sender;
        emit AddDomain(name, addr, domains[name].owner);
    }
}