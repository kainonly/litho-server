package van.api.entities

import javax.persistence.*

@Entity(name = "v_policy")
class Policy {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(columnDefinition = "int(10) unsigned")
    var id: Int? = 0



}