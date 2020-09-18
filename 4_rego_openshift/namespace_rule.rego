package regospike

get_all_fields(in) = fields {ex := {f | in[f]}; fields = ex}

get_field_with_default(in, field, defaults) = v {
  not in[field]  
  v := defaults[field]
} else = v {
  v := in[field]
}

get_object_with_defaults(in, defaults) = out {
  in_fields := {k | in[k]}
  def_fields := {k | defaults[k]}
  fields := in_fields | def_fields
  out := {f : v | fields[f]; v := get_field_with_default(in, f, defaults)}
}

labels_mutated(in) = out {
  out := get_object_with_defaults(in, {"ADDEDLABEL": "BINGO"})
}

metadata_field(in, field) = out {
  field == "annotations"
  out := labels_mutated(in[field])
} else = val {
  val := in[field]
}

metadata_mutated(in) = out {
  fields := get_all_fields(in)
  out := {f:v | fields[f]; v := metadata_field(in, f)}
}

ns_field(in, field) = out {
  field == "metadata"
  out := metadata_mutated(in[field])
} else = val {
  val := in[field]
}

ns_mutated(in) = out {
  fields := get_all_fields(in)
  out := {f:v | fields[f]; v := ns_field(in, f)}
}
