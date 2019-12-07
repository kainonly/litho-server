package van.api.entities

import org.hibernate.type.BooleanType
import org.hibernate.type.TextType
import javax.persistence.*

@Entity(name = "v_acl")
class Acl() {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(columnDefinition = "int(10) unsigned")
    var id: Int? = 0

    @Column(name = "`key`", unique = true, nullable = false, length = 40)
    var key: String? = null

    @Column(nullable = false)
    var name: TextType? = null

    @Column(name = "`write`", nullable = false)
    var write: TextType? = null

    @Column(name = "`read`", nullable = false)
    var read: TextType? = null

    @Column(nullable = false)
    var status: BooleanType? = null

    @Column(nullable = false, columnDefinition = "int(10) unsigned default 0")
    var createTime: Int = 0

    @Column(nullable = false, columnDefinition = "int(10) unsigned default 0")
    var updateTime: Int = 0
}

