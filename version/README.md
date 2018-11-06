Structure of the versiontable.yml file

versiontable:
  {FEATURE KEYWORD}:
    startversion: {Start version}
    description: {Description about feature}


For example, the file structure can be:

versiontable:
  CREATED_VERSION_MANAGER:
    startversion: v0.1.1
    description: Create version manager for Insolar platform 
  CHANGED_CONSENSUS_VERSION_MANAGER:
      startversion: v0.5.2
      description: Changed consensus of the version manager to BFT
...
