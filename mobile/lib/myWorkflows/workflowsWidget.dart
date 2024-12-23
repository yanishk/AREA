import 'package:flutter/material.dart';

class WorkflowCard extends StatelessWidget {
  const WorkflowCard({
    super.key,
    required this.content,
    this.color = Colors.black,
    this.id = "",
  });

  final Color color;
  final String id;
  final Widget content;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 15, vertical: 10),
      child: LayoutBuilder(
        builder: (context, constraints) {
          return Container(
            decoration: BoxDecoration(
              border: Border.all(
                color: color,
                width: 2.0,
              ),
              borderRadius: BorderRadius.circular(8.0),
            ),
            child: Column(
              children: [
                DecoratedBox(
                  decoration: BoxDecoration(
                    color: color,
                  ),
                  child: Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 25, vertical: 10),
                    child: Padding(
                      padding: EdgeInsets.only(top: 15),
                      child: SizedBox(
                        width: constraints.maxWidth,
                        child: content,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          );
        },
      ),
    );
  }
}