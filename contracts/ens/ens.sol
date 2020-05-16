pragma solidity 0.5.15;

contract Test {
    struct Domain {
        address pointTo;
        address owner;
    }

    mapping(bytes32 => Domain) public domains;

    function alter(bytes32 memory name, address addr) public {
        
    }
}