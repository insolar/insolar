package core

// ArtifactManager is a high level storage interface.
type ArtifactManager interface {
	// SetArchPref stores a list of preferred VM architectures memory.
	//
	// When returning classes storage will return compiled code according to this preferences. VM is responsible for
	// calling this method before fetching object in a new process. If preference is not provided, object getters will
	// return an error.
	SetArchPref(pref []MachineType)

	// GetExactObj returns code and memory of provided object/class state. Deactivation records should be ignored
	// (e.g. object considered to be active).
	//
	// This method is used by validator to fetch the exact state of the object that was used by the executor.
	GetExactObj(class, object RecordRef) ([]byte, []byte, error)

	// GetLatestObj returns descriptors for latest known state of the object/class known to the storage. The caller
	// should provide latest known states of the object/class known to it. If the object or the class is deactivated,
	// an error should be returned.
	//
	// Returned descriptors will provide methods for fetching migrations and appends relative to the provided states.
	GetLatestObj(object, storedClassState, storedObjState RecordRef) (
		ClassDescriptor, ObjectDescriptor, error,
	)

	// DeployCode creates new code record in storage.
	//
	// Code records are used to activate class or as migration code for an object.
	DeployCode(domain, request RecordRef, codeMap map[MachineType][]byte) (RecordRef, error)

	// ActivateClass creates activate class record in storage. Provided code reference will be used as a class code
	// and memory as the default memory for class objects.
	//
	// Activation reference will be this class'es identifier and referred as "class head".
	ActivateClass(domain, request, code RecordRef, memory []byte) (RecordRef, error)

	// DeactivateClass creates deactivate record in storage. Provided reference should be a reference to the head of
	// the class. If class is already deactivated, an error should be returned.
	//
	// Deactivated class cannot be changed or instantiate objects.
	DeactivateClass(domain, request, class RecordRef) (RecordRef, error)

	// UpdateClass creates amend class record in storage. Provided reference should be a reference to the head of
	// the class. Migrations are references to code records.
	//
	// Migration code will be executed by VM to migrate objects memory in the order they appear in provided slice.
	UpdateClass(domain, request, class, code RecordRef, migrationRefs []RecordRef) (
		RecordRef, error,
	)

	// ActivateObj creates activate object record in storage. Provided class reference will be used as objects class
	// memory as memory of crated object. If memory is not provided, the class default memory will be used.
	//
	// Activation reference will be this object's identifier and referred as "object head".
	ActivateObj(domain, request, class RecordRef, memory []byte) (RecordRef, error)

	// DeactivateObj creates deactivate object record in storage. Provided reference should be a reference to the head
	// of the object. If object is already deactivated, an error should be returned.
	//
	// Deactivated object cannot be changed.
	DeactivateObj(domain, request, obj RecordRef) (RecordRef, error)

	// UpdateObj creates amend object record in storage. Provided reference should be a reference to the head of the
	// object. Provided memory well be the new object memory.
	//
	// This will nullify all the object's append delegates. VM is responsible for collecting all appends and adding
	// them to the new memory manually if its required.
	UpdateObj(domain, request, obj RecordRef, memory []byte) (RecordRef, error)

	// AppendObjDelegate creates append object record in storage. Provided reference should be a reference to the head
	// of the object. Provided memory well be used as append delegate memory.
	//
	// Object's delegates will be provided by GetLatestObj. Any object update will nullify all the object's append
	// delegates. VM is responsible for collecting all appends and adding them to the new memory manually if its
	// required.
	AppendObjDelegate(domain, request, obj RecordRef, memory []byte) (RecordRef, error)
}

type ClassDescriptor interface {
	// GetCode fetches the latest class code known to storage. Code will be fetched according to architecture preferences
	// set via SetArchPref in artifact manager. If preferences are not provided, an error will be returned.
	GetCode() ([]byte, error)

	// GetMigrations fetches all migrations from provided to artifact manager state to the last state known to storage. VM
	// is responsible for applying these migrations and updating objects.
	GetMigrations() ([][]byte, error)
}

type ObjectDescriptor interface {
	// GetMemory fetches latest memory of the object known to storage.
	GetMemory() ([]byte, error)
	// GetDelegates fetches unamended delegates from storage.
	//
	// VM is responsible for collecting all delegates and adding them to the object memory manually if its required.
	GetDelegates() ([][]byte, error)
}
