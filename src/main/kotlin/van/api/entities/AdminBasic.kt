package van.api.entities

import javax.persistence.*

@Entity(name = "v_admin_basic")
class AdminBasic {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(columnDefinition = "int(10) unsigned")
    var id: Int? = 0


}